package cdkey

import (
	pkgBean "webapi/bean"
)

func GenerateInviteCode(channel string, num int, deadline int64) (error, []string) {
	results := []string{}

	for num > len(results) {
		cdkeys := generateCDKeyString(num-len(results), 8)
		err, batch := pkgBean.BatchAddInviteCode(channel, deadline, cdkeys)
		if err != nil {
			return err, nil
		}

		results = append(results, batch...)
	}

	return nil, results
}
