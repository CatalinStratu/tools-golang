// Package saver2v2 contains functions to render and write a tag-value
// formatted version of an in-memory SPDX document and its sections
// (version 2.2).
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package saver2v2

import (
	"fmt"
	"io"
	"sort"

	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2_2"
)

// RenderDocument2_2 is the main entry point to take an SPDX in-memory
// Document (version 2.2), and render it to the received io.Writer.
// It is only exported in order to be available to the tvsaver package,
// and typically does not need to be called by client code.
func RenderDocument2_2(doc *v2_2.Document, w io.Writer) error {
	if doc.CreationInfo == nil {
		return fmt.Errorf("document had nil CreationInfo section")
	}

	if doc.SPDXVersion != "" {
		fmt.Fprintf(w, "SPDXVersion: %s\n", doc.SPDXVersion)
	}
	if doc.DataLicense != "" {
		fmt.Fprintf(w, "DataLicense: %s\n", doc.DataLicense)
	}
	if doc.SPDXIdentifier != "" {
		fmt.Fprintf(w, "SPDXID: %s\n", common.RenderElementID(doc.SPDXIdentifier))
	}
	if doc.DocumentName != "" {
		fmt.Fprintf(w, "DocumentName: %s\n", doc.DocumentName)
	}
	if doc.DocumentNamespace != "" {
		fmt.Fprintf(w, "DocumentNamespace: %s\n", doc.DocumentNamespace)
	}
	// print EDRs in order sorted by identifier
	sort.Slice(doc.ExternalDocumentReferences, func(i, j int) bool {
		return doc.ExternalDocumentReferences[i].DocumentRefID < doc.ExternalDocumentReferences[j].DocumentRefID
	})
	for _, edr := range doc.ExternalDocumentReferences {
		fmt.Fprintf(w, "ExternalDocumentRef: DocumentRef-%s %s %s:%s\n",
			edr.DocumentRefID, edr.URI, edr.Checksum.Algorithm, edr.Checksum.Value)
	}
	if doc.DocumentComment != "" {
		fmt.Fprintf(w, "DocumentComment: %s\n", textify(doc.DocumentComment))
	}

	renderCreationInfo2_2(doc.CreationInfo, w)

	if len(doc.Files) > 0 {
		fmt.Fprintf(w, "##### Unpackaged files\n\n")
		sort.Slice(doc.Files, func(i, j int) bool {
			return doc.Files[i].FileSPDXIdentifier < doc.Files[j].FileSPDXIdentifier
		})
		for _, fi := range doc.Files {
			renderFile2_2(fi, w)
		}
	}

	// sort Packages by identifier
	sort.Slice(doc.Packages, func(i, j int) bool {
		return doc.Packages[i].PackageSPDXIdentifier < doc.Packages[j].PackageSPDXIdentifier
	})
	for _, pkg := range doc.Packages {
		fmt.Fprintf(w, "##### Package: %s\n\n", pkg.PackageName)
		renderPackage2_2(pkg, w)
	}

	if len(doc.OtherLicenses) > 0 {
		fmt.Fprintf(w, "##### Other Licenses\n\n")
		for _, ol := range doc.OtherLicenses {
			renderOtherLicense2_2(ol, w)
		}
	}

	if len(doc.Relationships) > 0 {
		fmt.Fprintf(w, "##### Relationships\n\n")
		for _, rln := range doc.Relationships {
			renderRelationship2_2(rln, w)
		}
		fmt.Fprintf(w, "\n")
	}

	if len(doc.Annotations) > 0 {
		fmt.Fprintf(w, "##### Annotations\n\n")
		for _, ann := range doc.Annotations {
			renderAnnotation2_2(ann, w)
			fmt.Fprintf(w, "\n")
		}
	}

	if len(doc.Reviews) > 0 {
		fmt.Fprintf(w, "##### Reviews\n\n")
		for _, rev := range doc.Reviews {
			renderReview2_2(rev, w)
		}
	}

	return nil
}
