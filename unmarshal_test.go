package xmp

import (
	"testing"
)

func TestUnmarshal(t *testing.T) {
	packet, err := Unmarshal([]byte(xmpExample))
	if err != nil {
		t.Error(err.Error())
		return
	}
	if packet == nil {
		t.Error("packet can't be nil")
	}
}

const xmpExample = `<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/" x:xmptk="Adobe XMP Core 5.4-c005 78.147326, 2012/08/23-13:03:03        ">
   <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
      <rdf:Description rdf:about=""
            xmlns:dc="http://purl.org/dc/elements/1.1/"
            xmlns:xmp="http://ns.adobe.com/xap/1.0/"
            xmlns:xmpMM="http://ns.adobe.com/xap/1.0/mm/"
            xmlns:stRef="http://ns.adobe.com/xap/1.0/sType/ResourceRef#"
            xmlns:stEvt="http://ns.adobe.com/xap/1.0/sType/ResourceEvent#"
            xmlns:pdf="http://ns.adobe.com/pdf/1.3/">
         <dc:format>application/pdf</dc:format>
         <dc:title>
            <rdf:Alt>
               <rdf:li xml:lang="x-default">What Climate Change Means for Massachusetts</rdf:li>
            </rdf:Alt>
         </dc:title>
         <dc:creator>
            <rdf:Seq>
               <rdf:li>US EPA</rdf:li>
               <rdf:li>OAR</rdf:li>
               <rdf:li>Climate Change Division</rdf:li>
            </rdf:Seq>
         </dc:creator>
         <dc:description>
            <rdf:Alt>
               <rdf:li xml:lang="x-default">This fact sheet provides a concise overview of the observed and projected effects and impacts of climate change on Massachusetts.</rdf:li>
            </rdf:Alt>
         </dc:description>
         <dc:subject>
            <rdf:Bag>
               <rdf:li>EPA</rdf:li>
               <rdf:li>climate change</rdf:li>
               <rdf:li>state</rdf:li>
               <rdf:li>impacts</rdf:li>
               <rdf:li>fact sheet</rdf:li>
               <rdf:li>summary</rdf:li>
            </rdf:Bag>
         </dc:subject>
         <xmp:CreateDate>2016-04-03T14:29:11-04:00</xmp:CreateDate>
         <xmp:MetadataDate>2016-08-24T11:03:55-04:00</xmp:MetadataDate>
         <xmp:ModifyDate>2016-08-24T11:03:55-04:00</xmp:ModifyDate>
         <xmp:CreatorTool>Adobe InDesign CC 2015 (Windows)</xmp:CreatorTool>
         <xmpMM:InstanceID>uuid:ab12ab15-97ee-4dbc-99d4-d286053991dc</xmpMM:InstanceID>
         <xmpMM:OriginalDocumentID>xmp.did:8dc6091d-ba7e-6d4a-89c2-d58847cb2d43</xmpMM:OriginalDocumentID>
         <xmpMM:DocumentID>xmp.id:947d03c1-f44e-354a-a341-4fb8d9199775</xmpMM:DocumentID>
         <xmpMM:RenditionClass>proof:pdf</xmpMM:RenditionClass>
         <xmpMM:DerivedFrom rdf:parseType="Resource">
            <stRef:instanceID>xmp.iid:a9d0dc65-36c9-be45-a66c-672ce9ee6a45</stRef:instanceID>
            <stRef:documentID>xmp.did:8e25c5db-4978-3d4a-84bc-324a454ea383</stRef:documentID>
            <stRef:originalDocumentID>xmp.did:8dc6091d-ba7e-6d4a-89c2-d58847cb2d43</stRef:originalDocumentID>
            <stRef:renditionClass>default</stRef:renditionClass>
         </xmpMM:DerivedFrom>
         <xmpMM:History>
            <rdf:Seq>
               <rdf:li rdf:parseType="Resource">
                  <stEvt:action>converted</stEvt:action>
                  <stEvt:parameters>from application/x-indesign to application/pdf</stEvt:parameters>
                  <stEvt:softwareAgent>Adobe InDesign CC 2015 (Windows)</stEvt:softwareAgent>
                  <stEvt:changed>/</stEvt:changed>
                  <stEvt:when>2016-04-03T14:29:11-04:00</stEvt:when>
               </rdf:li>
            </rdf:Seq>
         </xmpMM:History>
         <pdf:Producer>Adobe PDF Library 15.0</pdf:Producer>
         <pdf:Trapped>False</pdf:Trapped>
         <pdf:Keywords>EPA; climate change; state; impacts; fact sheet; summary</pdf:Keywords>
      </rdf:Description>
   </rdf:RDF>
</x:xmpmeta>





















<?xpacket end="w"?>`
