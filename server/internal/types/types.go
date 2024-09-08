package types

import (
	"math/rand"
	"time"

	faker "github.com/brianvoe/gofakeit"
)

type Examination struct {
	Sender          string
	ExaminationType string
	PatientsName    string
	IsSenderADoctor bool
}

type ExaminationType string

const (
	Hip   string = "hip"
	Knee  string = "knee"
	Elbow string = "elbow"

	ExchangeName    string = "examinations"
	LogExchangeName string = "logs"
)

func init() {
	faker.Seed(0)
}

func RandomExamination() Examination {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	name := faker.Name()

	switch r.Intn(3) {
	case 0:
		return Examination{ExaminationType: Hip, PatientsName: name, IsSenderADoctor: true}
	case 1:
		return Examination{ExaminationType: Knee, PatientsName: name, IsSenderADoctor: true}
	case 2:
		return Examination{ExaminationType: Elbow, PatientsName: name, IsSenderADoctor: true}
	default:
		return Examination{ExaminationType: Hip, PatientsName: name, IsSenderADoctor: true}
	}
}
