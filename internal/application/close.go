package application

type CloseFunc func() error

var closes []CloseFunc

func Close() {
	for i := 0; i < len(closes); i++ {
		err := closes[i]()
		if err != nil {
			log.Err(err).Send()
		}
	}
}
