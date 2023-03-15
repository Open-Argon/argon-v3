package main

func Map(val anymap) ArObject {
	return ArObject{
		TYPE: "map",
		obj:  val,
	}
}
