package iptc

import (
	"github.com/adampresley/adamgokit/slices"
	"github.com/adampresley/imagemetadata/imagemodel"
)

var standardTags = map[string]func(imageData *imagemodel.ImageData, data []string){
	"2:25": func(imageData *imagemodel.ImageData, data []string) {
		if len(data) > 0 {
			imageData.Keywords = slices.Merge(imageData.Keywords, data)
		}
	},
	"2:80": func(imageData *imagemodel.ImageData, data []string) {
		if imageData.Author == "" && len(data) > 0 {
			imageData.Author = data[0]
		}
	},
	"2:100": func(imageData *imagemodel.ImageData, data []string) {
		if imageData.CountryCode == "" && len(data) > 0 {
			imageData.CountryCode = data[0]
		}
	},
	"2:101": func(imageData *imagemodel.ImageData, data []string) {
		if imageData.CountryName == "" && len(data) > 0 {
			imageData.CountryName = data[0]
		}
	},
	"2:105": func(imageData *imagemodel.ImageData, data []string) {
		if imageData.Headline == "" && len(data) > 0 {
			imageData.Headline = data[0]
		}
	},
	"2:90": func(imageData *imagemodel.ImageData, data []string) {
		if imageData.City == "" && len(data) > 0 {
			imageData.City = data[0]
		}
	},
	"2:95": func(imageData *imagemodel.ImageData, data []string) {
		if imageData.StateProvince == "" && len(data) > 0 {
			imageData.StateProvince = data[0]
		}
	},
	"2:92": func(imageData *imagemodel.ImageData, data []string) {
		if imageData.Address == "" && len(data) > 0 {
			imageData.Address = data[0]
		}
	},
	"2:120": func(imageData *imagemodel.ImageData, data []string) {
		if imageData.CaptionIPTC == "" && len(data) > 0 {
			imageData.CaptionIPTC = data[0]
		}
	},
	"2:116": func(imageData *imagemodel.ImageData, data []string) {
		if imageData.Copyright == "" && len(data) > 0 {
			imageData.Copyright = data[0]
		}
	},
	"2:5": func(imageData *imagemodel.ImageData, data []string) {
		if imageData.TitleIPTC == "" && len(data) > 0 {
			imageData.TitleIPTC = data[0]
		}
	},
}
