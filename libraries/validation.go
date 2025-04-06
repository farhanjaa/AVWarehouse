package libraries

import (
	"database/sql"
	"reflect"
	"strings"

	"github.com/farhanjaa/AVWarehouse/config"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type validation struct {
	conn *sql.DB
}

func NewValidation() *validation {
	conn, err := config.ConnectDB()

	if err != nil {
		panic(err)
	}

	return &validation{
		conn: conn,
	}

}

func (v *validation) Init() (*validator.Validate, ut.Translator) {

	//memanggil package translator
	translator := en.New()
	uni := ut.New(translator, translator)

	trans, _ := uni.GetTranslator("en")

	validate := validator.New()

	// register default translator (en)
	en_translations.RegisterDefaultTranslations(validate, trans)

	// mengubah label default
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		labelName := field.Tag.Get("label")
		return labelName
	})

	// Mengubah label default
	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} tidak boleh kosong", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t

	})

	validate.RegisterValidation("isunique", func(fl validator.FieldLevel) bool {
		params := fl.Param()
		split_params := strings.Split(params, "-")

		tableName := split_params[0]
		tableField := split_params[1]
		fieldValue := fl.Field().String()

		return v.checkIsUnique(tableName, tableField, fieldValue)
	})
	validate.RegisterTranslation("isunique", trans, func(ut ut.Translator) error {
		return ut.Add("isunique", "{0} sudah digunakan", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("isunique", fe.Field())
		return t

	})

	return validate, trans
}

func (v *validation) Struct(s interface{}) interface{} {

	validate, trans := v.Init()

	vErrors := make(map[string]interface{})

	err := validate.Struct(s)

	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			vErrors[e.StructField()] = e.Translate(trans)
		}
	}

	if len(vErrors) > 0 {
		return vErrors
	}

	return nil
}

func (v *validation) checkIsUnique(tableName, fieldName, fieldValue string) bool {

	row, _ := v.conn.Query("Select "+fieldName+" from "+tableName+" where "+fieldName+" = ?", fieldValue)

	defer row.Close()

	var result string
	for row.Next() {
		row.Scan(&result)
	}

	// email@tentangkode
	//email@tentangkode

	return result != fieldValue
}
