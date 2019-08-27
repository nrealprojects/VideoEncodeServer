package main

import "projects/VideoEncodeServer/models"

// BinarySearch : binary search,
func BinarySearch(arr []models.SlamPose, low int, high int, key uint64) int {
	mid := (low + high) / 2
	resultIndex := -1
	if low > high || key < arr[0].TimeNanos || key > arr[len(arr)-1].TimeNanos {
		resultIndex = -1
	} else {
		result := 0
		flag := false
		if flag, result = CheckValue(arr, mid, key); flag {
			resultIndex = result
		} else if arr[mid].TimeNanos > key {
			resultIndex = BinarySearch(arr, low, mid-1, key)
		} else {
			resultIndex = BinarySearch(arr, mid+1, high, key)
		}
	}

	return resultIndex
}

// CheckValue : check value is equal
func CheckValue(arr []models.SlamPose, index int, key uint64) (bool, int) {
	flag := false
	resultIndex := index
	if arr[index].TimeNanos == key {
		resultIndex = index
		flag = true
	} else if arr[index].TimeNanos < key && (index+1) < len(arr) && arr[index+1].TimeNanos > key {
		if (key - arr[index].TimeNanos) < (arr[index+1].TimeNanos - key) {
			resultIndex = index
		} else {
			resultIndex = index + 1
		}
		flag = true
	} else if arr[index].TimeNanos > key && (index-1) >= 0 && arr[index-1].TimeNanos < key {
		if (arr[index].TimeNanos - key) < (key - arr[index-1].TimeNanos) {
			resultIndex = index
		} else {
			resultIndex = (index - 1)
		}
		flag = true
	}
	return flag, resultIndex
}
