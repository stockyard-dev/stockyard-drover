package server
type Limits struct{Tier string;Description string}
func LimitsFor(tier string)Limits{if tier=="pro"{return Limits{Tier:"pro"}};return Limits{Tier:"free"}}
func(l Limits)IsPro()bool{return l.Tier=="pro"}
