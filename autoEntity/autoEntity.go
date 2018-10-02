package autoEntity

import (
	"fmt"
	"regexp"
	"strings"
)

func underLine2Camel(field string, isTableName bool) string {
	s := field
	point := 0
	index := 0
	length := len(s)
	var points []int
	for index != -1 {
		index = strings.Index(s, "_")
		point += index
		points = append(points, point)
		//fmt.Printf("string<==>point<==>index<==>%s, %d, %d\n", s, point, index)
		if point == length-1 {
			break
		}
		s = s[index+1:]
		point += 1
	}
	//fmt.Println(points)
	characters := strings.Split(field, "")
	for _, element := range points {
		if element != length-1 && element != -1 {
			characters[element+1] = strings.ToUpper(characters[element+1])
			//fmt.Println(characters[element+1])
		}
	}
	if isTableName {
		characters[0] = strings.ToUpper(characters[0])
	}
	s = strings.Join(characters, "")
	s = strings.Replace(s, "_", "", -1)
	//fmt.Println(s)
	return s

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

func getJavaCode(table string, fT map[string]string, fC map[string]string, sql string) string {
	var code string
	code += "@Getter\r\n@Setter\r\n"
	code += "public class " + table + "VO" + " {\n"
	for key, value := range fT {
		code += "    // " + fC[key] + "\n"
		code += "    private " + value + " " + key + ";\n"
		code += "\n"
	}
	code += "\n}\n"
	code += "/**\n"
	code += sql
	code += "\n"
	code += "**/"
	return code
}

func Generate(sql string) string {
	//匹配表名
	tnre := regexp.MustCompile("\"\\w+\"\\.\"(?P<Table>\\w+)\"\\ *\\(")
	matchs := tnre.FindStringSubmatch(sql)
	tableName := underLine2Camel(matchs[1], true)

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
			fieldType[underLine2Camel(_field, false)] = _type
		} else if matched := comRe.FindStringSubmatch(line); len(matched) >= 3 {
			fieldCom[underLine2Camel(matched[1], false)] = matched[2]
		}
	}
	fmt.Printf("count ==>%d fieldMap ==>%v \n", len(fieldType), fieldType)
	fmt.Printf("count ==>%d commentMap ==>%v \n", len(fieldCom), fieldCom)

	code := getJavaCode(tableName, fieldType, fieldCom, sql)
	return code
}
