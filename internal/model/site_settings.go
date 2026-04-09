package model

type SiteSettingsRecord struct {
	ID             int64 `gorm:"primaryKey;autoIncrement:false"`
	AllowPDFExport bool  `gorm:"column:allow_pdf_export;not null;default:true"`
}

func (SiteSettingsRecord) TableName() string {
	return "site_settings"
}
