package ofxio

import (
	"encoding/xml"
	"io"
	"io/ioutil"

	"github.com/shopspring/decimal"
	. "github.com/stevegt/goadapt"
)

type Doc struct {
	XMLName xml.Name `xml:"OFX"`
	SONRS   SONRS    `xml:"SIGNONMSGSRSV1>SONRS"`
	// BANKMSGSRSV1       BANKMSGSRSV1       `xml:"BANKMSGSRSV1,omitempty"`
	CREDITCARDMSGSRSV1 CREDITCARDMSGSRSV1 `xml:"CREDITCARDMSGSRSV1,omitempty"`
}

type SONRS struct {
	CODE     int    `xml:"STATUS>CODE"`
	SEVERITY string `xml:"STATUS>SEVERITY"`
	DTSERVER string `xml:"DTSERVER"`
	LANGUAGE string `xml:"LANGUAGE"`
	DTPROFUP string `xml:"DTPROFUP"`
	ORG      string `xml:"FI>ORG"`
	FID      string `xml:"FI>FID"`
	INTUBID  string `xml:"INTU.BID,omitempty"`
}

type BANKMSGSRSV1 struct {
	CURDEF   string    `xml:"STMTTRNRS>STMTRS>CURDEF"`
	BANKID   string    `xml:"STMTTRNRS>STMTRS>BANKACCTFROM>BANKID"`
	ACCID    string    `xml:"STMTTRNRS>STMTRS>BANKACCTFROM>ACCTID"`
	ACCTTYPE string    `xml:"STMTTRNRS>STMTRS>BANKACCTFROM>ACCTTYPE"`
	Txs      []*BankTx `xml:"STMTTRNRS>STMTRS>BANKTRANLIST>STMTTRN"`
}

type CREDITCARDMSGSRSV1 struct {
	TRNUID      string          `xml:"CCSTMTTRNRS>TRNUID"`
	CODE        string          `xml:"CCSTMTTRNRS>STATUS>CODE"`
	SEVERITY    string          `xml:"CCSTMTTRNRS>STATUS>SEVERITY"`
	CURDEF      string          `xml:"CCSTMTTRNRS>CCSTMTRS>CURDEF"`
	ACCTID      string          `xml:"CCSTMTTRNRS>CCSTMTRS>CCACCTFROM>ACCTID"`
	DTSTART     string          `xml:"CCSTMTTRNRS>CCSTMTRS>BANKTRANLIST>DTSTART"`
	DTEND       string          `xml:"CCSTMTTRNRS>CCSTMTRS>BANKTRANLIST>DTEND"`
	Txs         []*CcTx         `xml:"CCSTMTTRNRS>CCSTMTRS>BANKTRANLIST>STMTTRN"`
	LEDGERBAL   decimal.Decimal `xml:"CCSTMTTRNRS>CCSTMTRS>LEDGERBAL>BALAMT"`
	LEDGERBALDT string          `xml:"CCSTMTTRNRS>CCSTMTRS>LEDGERBAL>DTASOF"`
	AVAILBAL    decimal.Decimal `xml:"CCSTMTTRNRS>CCSTMTRS>AVAILBAL>BALAMT"`
	AVAILBALDT  string          `xml:"CCSTMTTRNRS>CCSTMTRS>AVAILBAL>DTASOF"`
}

type BankTx struct {
	XMLName  xml.Name        `xml:"STMTTRN"`
	TRNTYPE  string          `xml:"TRNTYPE"`
	DTPOSTED string          `xml:"DTPOSTED"`
	DTUSER   string          `xml:"DTUSER"`
	TRNAMT   decimal.Decimal `xml:"TRNAMT"`
	FITID    string          `xml:"FITID"`
	NAME     string          `xml:"NAME"`
	MEMO     string          `xml:"MEMO"`
	BANKID   string          `xml:"BANKACCTTO>BANKID"`
	ACCTID   string          `xml:"BANKACCTTO>ACCTID"`
	ACCTTYPE string          `xml:"BANKACCTTO>ACCTTYPE"`
}

type CcTx struct {
	XMLName  xml.Name        `xml:"STMTTRN"`
	TRNTYPE  string          `xml:"TRNTYPE"`
	DTPOSTED string          `xml:"DTPOSTED"`
	DTUSER   string          `xml:"DTUSER"`
	TRNAMT   decimal.Decimal `xml:"TRNAMT"`
	FITID    string          `xml:"FITID"`
	NAME     string          `xml:"NAME"`
	ACCTID   string          `xml:"CCACCTTO>ACCTID"`
	MEMO     string          `xml:"MEMO"`
}

func Import(r io.Reader) (doc *Doc, err error) {
	defer Return(&err)

	buf, err := ioutil.ReadAll(r)
	Ck(err)

	doc = &Doc{}
	err = xml.Unmarshal(buf, doc)
	Ck(err)

	return doc, nil
}

func (doc *Doc) Export(w io.Writer) (err error) {
	defer Return(&err)

	header := `<?xml version="1.0" encoding="utf-8" ?>
<?OFX OFXHEADER="200" VERSION="202" SECURITY="NONE" OLDFILEUID="NONE" NEWFILEUID="NONE"?>
`
	_, err = w.Write([]byte(header))
	Ck(err)

	enc := xml.NewEncoder(w)
	enc.Indent("", "    ")

	err = enc.Encode(doc)
	Ck(err)

	return nil
}
