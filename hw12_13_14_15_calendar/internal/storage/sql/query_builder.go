package sqlstorage

import "strings"

type UpdateQueryBuilder struct {
	table       string
	setClause   []string
	whereClause []string
}

func NewUpdateQueryBuilder(table string) UpdateQueryBuilder {
	return UpdateQueryBuilder{
		table: table,
	}
}

func (q *UpdateQueryBuilder) Build() string {
	if len(q.setClause) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	sb.WriteString(q.table)
	sb.WriteString(" SET ")
	sb.WriteString(strings.Join(q.setClause, ", "))

	sb.WriteString(" WHERE ")
	sb.WriteString(strings.Join(q.whereClause, " "))

	return sb.String()
}

func (q *UpdateQueryBuilder) SetIf(condition bool, setClause string) {
	if condition {
		q.setClause = append(q.setClause, setClause)
	}
}

func (q *UpdateQueryBuilder) Where(whereClause string) {
	q.whereClause = append(q.whereClause, whereClause)
}
