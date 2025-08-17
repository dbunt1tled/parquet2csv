package schema

import (
	"github.com/dbunt1tled/parquet2csv/internal/helper"
)

type EnrichData struct {
	ID                                  int32  `parquet:"name=id, type=INT32"`
	PatientPharmacyID                   string `parquet:"name=patient_pharmacy_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	PharmacyID                          string `parquet:"name=pharmacy_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	CPatientDob                         string `parquet:"name=c_patient_dob, type=BYTE_ARRAY, convertedtype=UTF8"` // DATE
	PatientGender                       string `parquet:"name=patient_gender, type=BYTE_ARRAY, convertedtype=UTF8"`
	Race                                string `parquet:"name=race, type=BYTE_ARRAY, convertedtype=UTF8"`
	Ethnicity                           string `parquet:"name=ethnicity, type=BYTE_ARRAY, convertedtype=UTF8"`
	Guardian                            string `parquet:"name=guardian, type=BYTE_ARRAY, convertedtype=UTF8"`
	Address                             string `parquet:"name=address, type=BYTE_ARRAY, convertedtype=UTF8"`
	PhoneNumber                         string `parquet:"name=phone_number, type=BYTE_ARRAY, convertedtype=UTF8"`
	Email                               string `parquet:"name=email, type=BYTE_ARRAY, convertedtype=UTF"`
	MedName                             string `parquet:"name=med_name, type=BYTE_ARRAY, convertedtype=UTF8"`
	MedTA                               string `parquet:"name=med_ta, type=BYTE_ARRAY, convertedtype=UTF8"`
	MedDose                             string `parquet:"name=med_dose, type=BYTE_ARRAY, convertedtype=UTF8"`
	CMedGenericName                     string `parquet:"name=c_med_generic_name, type=BYTE_ARRAY, convertedtype=UTF8"`
	CMedStrength                        string `parquet:"name=c_med_strength, type=BYTE_ARRAY, convertedtype=UTF8"`
	Level1                              string `parquet:"name=level_1, type=BYTE_ARRAY, convertedtype=UTF8"`
	Level2                              string `parquet:"name=level_2, type=BYTE_ARRAY, convertedtype=UTF8"`
	CDispenseDate                       string `parquet:"name=c_dispense_date, type=BYTE_ARRAY, convertedtype=UTF8"` //nolint:lll    // DATE int32(time.Now().Unix() / 3600 / 24),
	MedRxNormID                         string `parquet:"name=med_rx_norm_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	MedNDCID                            string `parquet:"name=med_ndc_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	MedDescription                      string `parquet:"name=med_description, type=BYTE_ARRAY, convertedtype=UTF8"`
	DispenseQuantity                    string `parquet:"name=dispense_quantity, type=BYTE_ARRAY, convertedtype=UTF8"`
	DispenseQuantityUnit                string `parquet:"name=dispense_quantity_unit, type=BYTE_ARRAY, convertedtype=UTF8"`
	CDispenseDaysSupply                 string `parquet:"name=c_dispense_days_supply, type=BYTE_ARRAY, convertedtype=UTF8"` //nolint:lll    // INT32
	CDispenseQuantity                   string `parquet:"name=c_dispense_quantity, type=BYTE_ARRAY, convertedtype=UTF8"`
	PotentialPatientID                  int32  `parquet:"name=potential_patient_id, type=INT32"`                                                  //nolint:lll    // INT32
	PotentialPatientMedicationHistoryID string `parquet:"name=potential_patient_medication_history_id, type=BYTE_ARRAY, repetitiontype=OPTIONAL"` //nolint:lll    // INT32
	PharmacyNPI                         string `parquet:"name=pharmacy_npi, type=BYTE_ARRAY, convertedtype=UTF8"`
	PharmacyNCPDP                       string `parquet:"name=pharmacy_ncpdp, type=BYTE_ARRAY, convertedtype=UTF8"`
	PharmacyNABP                        string `parquet:"name=pharmacy_nabp, type=BYTE_ARRAY, convertedtype=UTF8"`
	ImportID                            string `parquet:"name=import_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	CMedCleanName                       string `parquet:"name=c_med_clean_name, type=BYTE_ARRAY, convertedtype=UTF8"`
	CDispenseQuantityUnit               string `parquet:"name=c_dispense_quantity_unit, type=BYTE_ARRAY, convertedtype=UTF8"` //nolint:lll
	Level3                              string `parquet:"name=level_3, type=BYTE_ARRAY, convertedtype=UTF8"`
	Level4                              string `parquet:"name=level_4, type=BYTE_ARRAY, convertedtype=UTF8"`
	CMedSourceCountry                   string `parquet:"name=c_med_source_country, type=BYTE_ARRAY, convertedtype=UTF8"`
	Level5                              string `parquet:"name=level_5, type=BYTE_ARRAY, convertedtype=UTF8"`
	OriginalMedName                     string `parquet:"name=original_med_name, type=BYTE_ARRAY, convertedtype=UTF8"`
	FileName                            string `parquet:"name=filename, type=BYTE_ARRAY, convertedtype=UTF8"`
	CreatedAt                           int32  `parquet:"name=created_at, type=INT32"` // INT32
	CompareKey                          string `parquet:"name=compare_key, type=BYTE_ARRAY, convertedtype=UTF8"`
	AdhDate                             string `parquet:"name=adh_date, type=BYTE_ARRAY, convertedtype=UTF8"` //nolint:lll    // DATE int32(time.Now().Unix() / 3600 / 24),
}

func ProcessEnrichData(record []string, sc interface{}, header []string) interface{} {
	eData := EnrichData{
		ID:                                  helper.StrToInt32(record[0], true),
		PatientPharmacyID:                   record[1],
		PharmacyID:                          record[2],
		CPatientDob:                         record[3],
		PatientGender:                       record[4],
		Race:                                record[5],
		Ethnicity:                           record[6],
		Guardian:                            record[7],
		Address:                             record[8],
		PhoneNumber:                         record[9],
		Email:                               record[10],
		MedName:                             record[11],
		MedTA:                               record[12],
		MedDose:                             record[13],
		CMedGenericName:                     record[14],
		CMedStrength:                        record[15],
		Level1:                              record[16],
		Level2:                              record[17],
		CDispenseDate:                       record[18],
		MedRxNormID:                         record[19],
		MedNDCID:                            record[20],
		MedDescription:                      record[21],
		DispenseQuantity:                    record[22],
		DispenseQuantityUnit:                record[23],
		CDispenseDaysSupply:                 record[24],
		PotentialPatientID:                  helper.StrToInt32(record[26], true),
		PotentialPatientMedicationHistoryID: record[27],
		PharmacyNPI:                         record[28],
		PharmacyNCPDP:                       record[29],
		PharmacyNABP:                        record[30],
		ImportID:                            record[31],
		CMedCleanName:                       record[32],
		CDispenseQuantityUnit:               record[33],
		Level3:                              record[34],
		Level4:                              record[35],
		CMedSourceCountry:                   record[36],
		Level5:                              record[37],
		OriginalMedName:                     record[38],
		FileName:                            record[39],
		CreatedAt:                           helper.StrToInt32(record[40], true),
		CompareKey:                          record[41],
		AdhDate:                             record[42],
	}

	return eData
}
