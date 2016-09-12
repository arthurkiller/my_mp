package conf

//make the conf struct
type Config struct {
	NewsServerAddr    string
	ProfileServerAddr string

	UIDInfo  []string
	UIDInfoS [][]string

	NewsIDInfo  []string
	NewsIDInfoS [][]string

	MeipaiIDInfo  []string
	MeipaiIDInfoS [][]string

	UIDBox  []string
	UIDBoxS [][]string

	UIDSelfbox  []string
	UIDSelfboxS [][]string

	UIDFans  []string
	UIDFansS [][]string

	UIDFollow  []string
	UIDFollowS [][]string

	Maxactive   int
	Maxidle     int
	Idletimeout int
}
