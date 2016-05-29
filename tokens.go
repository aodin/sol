package sol

// Tokens holds SQL tokens that are used for compilation / AST parsing
const (
	CROSSJOIN      = "CROSS JOIN"
	DELETE         = "DELETE"
	DISTINCT       = "DISTINCT"
	FROM           = "FROM"
	FULLOUTERJOIN  = "FULL OUTER JOIN"
	GROUPBY        = "GROUP BY"
	HAVING         = "HAVING"
	INNERJOIN      = "INNER JOIN"
	INSERT         = "INSERT"
	INTO           = "INTO"
	LEFTOUTERJOIN  = "LEFT OUTER JOIN"
	LIMIT          = "LIMIT"
	OFFSET         = "OFFSET"
	ORDERBY        = "ORDER BY"
	RIGHTOUTERJOIN = "RIGHT OUTER JOIN"
	SELECT         = "SELECT"
	SET            = "SET"
	UPDATE         = "UPDATE"
	VALUES         = "VALUES"
	WHERE          = "WHERE"
	WHITESPACE     = " "
)
