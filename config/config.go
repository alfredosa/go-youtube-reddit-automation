package config

import "time"

type Config struct {
	Goggle struct {
		API_Key       string
		CX            string
		Refresh_Token string
	}
	Directory struct {
		Path string
	}
	Background struct {
		Path string
	}
	YoutubeCredential struct {
		Client_ID     string
		Client_Secret string
	}
	AmazonAWSCredential struct {
		Aws_Access_Key_ID     string
		Aws_Secret_Access_Key string
		RegionName            string
	}
	Reddit struct {
		Subreddit     string
		Topn_Comments int
	}
	VideoSetup struct {
		Total_Video_Duration int
		Pause                float64
	}
	TextToSpeechSetup struct {
		Multiple_Voices bool
		Voice_ID        string
	}
	App struct {
		Upload_To_Youtube   bool
		Generation_Interval time.Duration
	}
	NewsAPI struct {
		API_Key string
	}
	Database struct {
		Port     int
		User     string
		Password string
		DB_Name  string
	}
}
