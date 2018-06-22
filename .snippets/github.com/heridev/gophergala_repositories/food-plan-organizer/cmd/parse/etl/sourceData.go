package etl

// Sources of Data File (file name = DATA_SRC).
//
// This file (Table 15) provides a citation to the DataSrc_ID in the Sources of Data Link file.
//
// - Links to Nutrient Data file by NDB No. through the Sources of Data Link file
//
// DataSrc_ID   A 6*    Unique number identifying the reference/source.
// Authors      A 255   List of authors for a journal article or name of sponsoring organization for other documents.
// Title        A 255   Title of article or name of document, such as a report from a company or trade association.
// Year         A 4     Year article or document was published.
// Journal      A 135   Name of the journal in which the article was published.
// Vol_City     A 16    Volume number for journal articles, books, or reports; city where sponsoring organization is located.
// Issue_State  A 5     Issue number for journal article; State where the sponsoring organization is located.
// Start_Page   A 5     Starting page number of article/document.
// End_Page     A 5     Ending page number of article/document.

type sourceData struct {
}
