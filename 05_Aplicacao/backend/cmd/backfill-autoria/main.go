package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const senateAuthoritiesURL = "https://legis.senado.leg.br/dadosabertos/senador/%s/autorias.json"

func main() {
	_ = godotenv.Load(".env")

	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Exec(`
		ALTER TABLE public.autoria
		ADD COLUMN IF NOT EXISTS principal BOOLEAN NOT NULL DEFAULT false
	`).Error; err != nil {
		log.Fatalf("erro ao adicionar coluna principal: %v", err)
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_autoria_autor_principal
		ON public.autoria (casa, id_autor, principal)
	`).Error; err != nil {
		log.Fatalf("erro ao criar índice de autoria: %v", err)
	}

	if err := db.Exec(`
		UPDATE public.autoria
		SET principal = (ordem_assinatura = 1)
		WHERE casa = 'CAMARA'
	`).Error; err != nil {
		log.Fatalf("erro ao atualizar autoria da Câmara: %v", err)
	}
	fmt.Println("Câmara atualizada (principal = ordem_assinatura 1).")

	var senatorRows []struct {
		IDAutor string
	}
	if err := db.Raw(`
		SELECT DISTINCT id_autor
		FROM public.autoria
		WHERE casa = 'SENADO'
		ORDER BY id_autor
	`).Scan(&senatorRows).Error; err != nil {
		log.Fatal(err)
	}

	client := &http.Client{Timeout: 60 * time.Second}

	for i, row := range senatorRows {
		senatorID := row.IDAutor
		codigo, err := strconv.Atoi(senatorID)
		if err != nil {
			fmt.Printf("[%d/%d] id inválido no Senado: %s\n", i+1, len(senatorRows), senatorID)
			continue
		}

		principais, err := buscarAutoriasPrincipais(client, codigo)
		if err != nil {
			fmt.Printf("[%d/%d] senador %s: erro na API: %v\n", i+1, len(senatorRows), senatorID, err)
			continue
		}

		if err := db.Exec(`
			UPDATE public.autoria
			SET principal = false
			WHERE casa = 'SENADO' AND id_autor = ?
		`, senatorID).Error; err != nil {
			log.Fatal(err)
		}

		for matterID := range principais {
			if err := db.Exec(`
				UPDATE public.autoria
				SET principal = true
				WHERE casa = 'SENADO' AND id_autor = ? AND id_proposicao = ?
			`, senatorID, matterID).Error; err != nil {
				log.Fatal(err)
			}
		}

		fmt.Printf("[%d/%d] senador %s: %d matéria(s) como autor principal\n",
			i+1, len(senatorRows), senatorID, len(principais))

		time.Sleep(150 * time.Millisecond)
	}

	var stats []struct {
		Casa      string
		Principal bool
		Total     int64
	}
	if err := db.Raw(`
		SELECT casa, principal, COUNT(*) AS total
		FROM public.autoria
		GROUP BY casa, principal
		ORDER BY casa, principal
	`).Scan(&stats).Error; err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nResumo final:")
	for _, s := range stats {
		fmt.Printf("  %s principal=%v: %d\n", s.Casa, s.Principal, s.Total)
	}
}

func buscarAutoriasPrincipais(client *http.Client, senatorID int) (map[string]bool, error) {
	url := fmt.Sprintf(senateAuthoritiesURL, strconv.Itoa(senatorID))

	var body []byte
	var lastErr error

	for tentativa := 1; tentativa <= 3; tentativa++ {
		resp, err := client.Get(url)
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(tentativa) * time.Second)
			continue
		}

		body, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(tentativa) * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("HTTP %d", resp.StatusCode)
			time.Sleep(time.Duration(tentativa) * time.Second)
			continue
		}

		lastErr = nil
		break
	}

	if lastErr != nil {
		return nil, lastErr
	}

	var payload any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	principais := map[string]bool{}
	coletarPrincipais(payload, principais)
	return principais, nil
}

func coletarPrincipais(obj any, principais map[string]bool) {
	switch v := obj.(type) {
	case map[string]any:
		materia, temMateria := v["Materia"].(map[string]any)
		indicador, temIndicador := v["IndicadorAutorPrincipal"].(string)
		if temMateria && temIndicador {
			codigo := fmt.Sprint(materia["Codigo"])
			if strings.EqualFold(strings.TrimSpace(indicador), "Sim") {
				principais[codigo] = true
			}
		}
		for _, filho := range v {
			coletarPrincipais(filho, principais)
		}
	case []any:
		for _, filho := range v {
			coletarPrincipais(filho, principais)
		}
	}
}
