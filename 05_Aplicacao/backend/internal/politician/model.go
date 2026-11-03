package politician

import "time"

type Politician struct {
	House           string  `gorm:"column:casa;primaryKey" json:"casa"`
	IDPolitician    string  `gorm:"column:id_parlamentar;primaryKey" json:"id_parlamentar"`
	ReferenceName   string  `gorm:"column:nome_referencia" json:"nome_referencia"`
	PhotoURL       *string `gorm:"column:url_foto" json:"url_foto"`
	Email          *string `gorm:"column:email" json:"email"`
}

func (Politician) TableName() string {
	return "public.parlamentar"
}

type House struct {
	House string `gorm:"column:casa;primaryKey" json:"casa"`
	Name  string `gorm:"column:nome" json:"nome"`
}

func (House) TableName() string {
	return "public.casa"
}

type State struct {
	UF   string `gorm:"column:uf;primaryKey" json:"uf"`
	Name string `gorm:"column:nome" json:"nome"`
}

func (State) TableName() string {
	return "public.uf"
}

type Party struct {
	Acronym string `gorm:"column:sigla;primaryKey" json:"sigla"`
}

func (Party) TableName() string {
	return "public.partido"
}

type Legislature struct {
	Number int       `gorm:"column:numero_legislatura;primaryKey" json:"numero_legislatura"`
	Start  time.Time `gorm:"column:inicio" json:"inicio"`
	End    time.Time `gorm:"column:fim" json:"fim"`
}

func (Legislature) TableName() string {
	return "public.legislatura"
}

type ObservedRegistration struct {
	House        string `gorm:"column:casa;primaryKey" json:"casa"`
	IDPolitician string `gorm:"column:id_parlamentar;primaryKey" json:"id_parlamentar"`
	IDLote       string `gorm:"column:id_lote;primaryKey" json:"id_lote"`
	Party        string `gorm:"column:partido" json:"partido"`
	UF           string `gorm:"column:uf" json:"uf"`
}

func (ObservedRegistration) TableName() string {
	return "public.cadastro_observado"
}

type PoliticianLegislature struct {
	House           string `gorm:"column:casa;primaryKey" json:"casa"`
	IDPolitician    string `gorm:"column:id_parlamentar;primaryKey" json:"id_parlamentar"`
	NumberLegislature int    `gorm:"column:numero_legislatura;primaryKey" json:"numero_legislatura"`
	InclusionCritL   string `gorm:"column:criterio_inclusao" json:"criterio_inclusao"`
	IDLote          string `gorm:"column:id_lote" json:"id_lote"`
}

func (PoliticianLegislature) TableName() string {
	return "public.parlamentar_legislatura"
}

type Term struct {
	IDTerm         string     `gorm:"column:id_exercicio" json:"id_exercicio"`
	IDPolitician   string     `gorm:"column:id_parlamentar" json:"id_parlamentar"`
	Start          time.Time  `gorm:"column:inicio" json:"inicio"`
	End            *time.Time `gorm:"column:fim" json:"fim"`
	CauseCode      string     `gorm:"column:codigo_causa" json:"codigo_causa"`
}

func (Term) TableName() string {
	return "public.exercicio"
}

