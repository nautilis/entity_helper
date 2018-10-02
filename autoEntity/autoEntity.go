package autoEntity

import (
	"fmt"
	"regexp"
	"strings"
)

type AutoEntity struct {
	fields    []string
	fieldType map[string]string
	fieldCom  map[string]string
	tableName string
	sql       string
}

func underLine2Camel(field string, isTableName bool) string {
	re := regexp.MustCompile("_")
	points := re.FindAllStringIndex(field, -1)
	fmt.Printf("string ==> %s, ponits ==>%v\n", field, points)
	length := len(field)
	characters := strings.Split(field, "")
	for _, point := range points {
		if point[1] != length {
			characters[point[1]] = strings.ToUpper(characters[point[1]])
		}
	}
	if isTableName {
		characters[0] = strings.ToUpper(characters[0])
	}
	res := strings.Join(characters, "")
	res = strings.Replace(res, "_", "", -1)
	return res
}

func getJavaType(sqlType string) string {
	var sqlTypeJavaType = make(map[string]string)
	sqlTypeJavaType["char"] = "String"
	sqlTypeJavaType["varchar"] = "String"
	sqlTypeJavaType["text"] = "String"
	sqlTypeJavaType["varying"] = "String"
	sqlTypeJavaType["timestamp"] = "Date"
	sqlTypeJavaType["time"] = "Date"
	sqlTypeJavaType["date"] = "Date"
	sqlTypeJavaType["bool"] = "Boolean"
	sqlTypeJavaType["boolean"] = "Boolean"

	ip := "numeric\\(\\d+\\)"
	bp := "numeric\\(\\d+,\\d+\\)"
	if matched, _ := regexp.MatchString(ip, sqlType); matched {
		return "Integer"
	} else if matched, _ := regexp.MatchString(bp, sqlType); matched {
		return "BigDecimal | Double"
	} else {
		_type := strings.Split(sqlType, "(")
		r := sqlTypeJavaType[_type[0]]
		if r != "" {
			return r
		}
		return sqlType
	}
}

func (a *AutoEntity) getJavaCode() string {
	var code string
	code += "@Getter\r\n@Setter\r\n"
	code += "public class " + a.tableName + "VO" + " {\n"
	for _, field := range a.fields {
		code += "    // " + a.fieldCom[field] + "\n"
		code += "    private " + a.fieldType[field] + " " + field + ";\n"
		code += "\n"
	}
	code += "\n}\n"
	code += "/**\n"
	code += a.sql
	code += "\n"
	code += "**/"
	return code
}

func Generate(sql string) string {
	//匹配表名
	tnre := regexp.MustCompile("\"\\w+\"\\.\"(?P<Table>\\w+)\"\\ *\\(")
	matchs := tnre.FindStringSubmatch(sql)
	tableName := underLine2Camel(matchs[1], true)

	var fields []string
	lines := strings.Split(sql, "\r\n")
	fieldType := make(map[string]string)
	fieldRe := "^\"\\w+\""
	fieldCom := make(map[string]string)
	//COMMENT ON COLUMN "adempiere"."lyy_income_sharing_set_log"."pre_rate" IS '收益分成比例  默认 -1 为未设置';
	comRe := regexp.MustCompile("COMMENT\\ ON\\ COLUMN\\ \"\\w+\"\\.\"\\w+\"\\.\"(?P<column>\\w+)\"\\ IS\\ '(?P<comment>(\\ *\\S+\\ *)+)'")
	for _, line := range lines {
		line = strings.Trim(line, "\r\n")
		line = strings.Trim(line, "\n")
		line = strings.Trim(line, " ")
		fmt.Println(line)
		if matched, _ := regexp.MatchString(fieldRe, line); matched {
			fieldList := strings.Split(line, " ")
			var _type string = getJavaType(fieldList[1])
			var _field string = strings.Trim(fieldList[0], "\"")
			_field = underLine2Camel(_field, false)
			fieldType[_field] = _type
			fields = append(fields, _field)
		} else if matched := comRe.FindStringSubmatch(line); len(matched) >= 3 {
			fieldCom[underLine2Camel(matched[1], false)] = matched[2]
		}
	}
	fmt.Printf("count ==>%d fieldMap ==>%v \n", len(fieldType), fieldType)
	fmt.Printf("count ==>%d commentMap ==>%v \n", len(fieldCom), fieldCom)

	entity := &AutoEntity{
		tableName: tableName,
		fieldType: fieldType,
		fieldCom:  fieldCom,
		sql:       sql,
		fields:    fields,
	}

	code := entity.getJavaCode()
	return code
}
