package client

func (b *bugyoClient) Punchmark() error {
	// move PunchMark page
	// GET: https://hromssp.obc.jp/{tenantCode}/{userCode}/timeclock/punchmark/

	// ClockIn
	// POST https://hromssp.obc.jp/{tenantCode}/{userCode}/TimeClock/InsertReadDateTime/
	//headers={'Referer': 'https://hromssp.obc.jp/{tenantCode}/{userCode}/timeclock/punchmark/', '__RequestVerificationToken': '<token>', 'X-Requested-With': 'XMLHttpRequest'}
	//data={'ClockType': 'ClockIn', 'LaborSystemID': '0', 'LaborSystemCode': '', 'LaborSystemName': '', 'PositionLatitude': <CompanyLatitude>, 'PositionLongitude': <CompanyLongitude>, 'PositionAccuracy': '0'}

	return nil
}
