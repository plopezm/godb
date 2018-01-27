package godb

// QueryGroup is a group of queries executed in one call
// to the database server, using multiple result sets.
type QueryGroup struct {
	db      *DB
	queries []*QueryLater
	results []*QueryResult
}

// QueryResult will containts data about a query after the queries
// of the QueryGroup are executed.
type QueryResult struct {
	err error
}

//
func (qr *QueryResult) Err() error {
	return qr.err
}

// QueryLater contains all informations to executes a query with a QueryGroup.
type QueryLater struct {
	err               error
	sql               string
	arguments         []interface{}
	recordDescription *recordDescription
	pointersGetter    pointersGetter
}

//
func (ql *QueryLater) Err() error {
	return ql.err
}

// NewQueryGroup initialize a new QueryGroup.
func (db *DB) NewQueryGroup() *QueryGroup {
	return &QueryGroup{
		db: db,
	}
}

// Add adds a QueryLater to a QueryGroup.
func (qg *QueryGroup) Add(q *QueryLater) *QueryResult {
	qg.queries = append(qg.queries, q)
	result := &QueryResult{}
	qg.results = append(qg.results, result)
	return result
}

// Do runs all que queryes in the QueryGroup, and fill all the QueryLater
// returned by Add with relared informations.
func (qg *QueryGroup) Do() error {
	// Merge all queries
	sqlLength := 2 * len(qg.queries) // ;\n
	argsCount := 0
	for _, ql := range qg.queries {
		sqlLength += len(ql.sql)
		argsCount += len(ql.arguments)
	}
	oneQueriesToRulesThemAll := NewSQLBuffer(sqlLength, argsCount)
	allRecords := make([]*recordDescription, 0, len(qg.queries))
	allPointersGetter := make([]pointersGetter, 0, len(qg.queries))
	for _, ql := range qg.queries {
		oneQueriesToRulesThemAll.Write(ql.sql, ql.arguments...)
		oneQueriesToRulesThemAll.Write(";\n")
		allRecords = append(allRecords, ql.recordDescription)
		allPointersGetter = append(allPointersGetter, ql.pointersGetter)
	}

	qg.db.doMultipleSelectOrWithReturning(
		oneQueriesToRulesThemAll.SQL(),
		oneQueriesToRulesThemAll.Arguments(),
		allRecords,
		allPointersGetter,
	)
	// run them
	// get results
	// fill all QueryResult
	//return fmt.Errorf("Not implemented... work in progress ;)")
	return nil
}
