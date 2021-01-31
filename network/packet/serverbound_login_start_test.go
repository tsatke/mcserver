package packet

func (suite *PacketSuite) TestServerboundLoginStart_Validate() {
	for _, test := range []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			"empty",
			"",
			true,
		},
		{
			"valid",
			"myUsername",
			false,
		},
		{
			"valid max len",
			"0123456789abcdef",
			false,
		},
		{
			"too long",
			"0123456789abcdefX",
			true,
		},
		{
			"valid utf8",
			"\u1000\u1001\u1002\u1003\u1004\u1005\u1006\u1007\u1008\u1009\u100a\u100b\u100c\u100d\u100e\u100f",
			false,
		},
		{
			"utf8 too long",
			"\u1000\u1001\u1002\u1003\u1004\u1005\u1006\u1007\u1008\u1009\u100a\u100b\u100c\u100d\u100e\u100f\u1234",
			true,
		},
	} {
		suite.Run(test.name, func() {
			err := ServerboundLoginStart{
				Username: test.username,
			}.Validate()
			if test.wantErr {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
		})
	}
}
