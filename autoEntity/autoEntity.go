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
	fmt.Printf("table name is %s\n", tableName)

	//"timestamp" varchar(12) COLLATE "default" NOT NULL,
	braceIndex := strings.Index(sql, "(")
	fmt.Println("braceIndex is %s \n", braceIndex)
	subSql := sql[braceIndex+2:]
	subList := strings.Split(subSql, ",\r\n")
	fieldsList := subList[:len(subList)-1]
	comment := subList[len(subList)-1]
	fmt.Printf("subSql <==> fieldsList <===> comment %s\n<====>%v\n <===> %s\n", subSql, fieldsList, comment)

	//处理字段
	fieldType := make(map[string]string)
	fieldRe := "^\"\\w+\""
	for _, ele := range fieldsList {
		ele = strings.Trim(ele, "\r\n")
		ele = strings.Trim(ele, "\n")
		ele = strings.Trim(ele, " ")
		fmt.Println(ele)
		if matched, _ := regexp.MatchString(fieldRe, ele); !matched {
			continue
		}
		eleList := strings.Split(ele, " ")
		var _type string = getJavaType(eleList[1])
		eleList[0] = strings.Trim(eleList[0], "\"")
		fieldType[underLine2Camel(eleList[0], false)] = _type
	}
	fmt.Printf("count <==> map, %d <==> %v\n", len(fieldType), fieldType)

	//处理注释
	comList := strings.Split(comment, ";")
	//COMMENT ON COLUMN "adempiere"."es_order"."actual_amount" IS '实际金额';
	comRe := regexp.MustCompile("COMMENT\\ ON\\ COLUMN\\ \"\\w+\"\\.\"\\w+\"\\.\"(?P<column>\\w+)\"\\ IS\\ '(?P<comment>\\S+)'")
	fieldCom := make(map[string]string)
	for _, com := range comList {
		com = strings.Trim(com, " ")
		fmt.Println(com)
		matched := comRe.FindStringSubmatch(com)
		fmt.Println(matched)
		if len(matched) == 3 {
			fieldCom[underLine2Camel(matched[1], false)] = matched[2]
		}
	}
	fmt.Printf("count <==> map, %d <==> %v \n", len(fieldCom), fieldCom)

	//拼接java code
	code := getJavaCode(tableName, fieldType, fieldCom, sql)
	return code

}
