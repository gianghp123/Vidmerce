package utils

import "github.com/jinzhu/copier"

func MapToDTO(source interface{}, dest interface{}) error {
	return copier.Copy(dest, source)
}

func MapToDTOs(source interface{}, dest interface{}) error {
	return copier.Copy(dest, source)
}
