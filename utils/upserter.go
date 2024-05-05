package utils

//
//type Upserter struct {
//	Table    string   `json:"table"`
//	Columns  []string `json:"columns"`
//	Conflict []string `json:"conflict"`
//	Upserts  []string `json:"upserts"`
//}
//
//func conflict(fields []string) string {
//	return fmt.Sprint(`ON CONFLICT ("`, strings.Join(fields, `","`), `") DO`)
//}
//
//func (u *Upserter) Build() string {
//	fmt.Sprintf(`INSERT INTO "%s" ("%s") VALUES `, u.Table, strings.Join(u.Columns, `","`))
//}
