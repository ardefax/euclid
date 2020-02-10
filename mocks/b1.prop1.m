(* ::Package:: *)

(* ::Input:: *)
(*coordA = {-1, 0}*)


(* ::Input:: *)
(*coordB = {1, 0}*)


(* ::Input:: *)
(*lineAB = Line[{coordA, coordB}]*)


(* ::Input:: *)
(*radius = EuclideanDistance[coordA, coordB]*)


(* ::Input:: *)
(*circleACE = Circle[coordA, radius]*)


(* ::Input:: *)
(*circleBCD = Circle[coordB, radius]*)


(* ::Input:: *)
(*pointA = Point[coordA]*)


(* ::Input:: *)
(*pointB= Point[coordB]*)


(* ::Input:: *)
(*intersections = {x,y} /. Solve[{x,y} \[Element] circleACE && {x,y} \[Element] circleBCD, {x,y}]*)


(* ::Input:: *)
(*coordC = intersections[[2]]*)


(* ::Input:: *)
(*lineAC = Line[{coordA, coordC}]*)


(* ::Input:: *)
(*lineBC = Line[{coordB, coordC}]*)


(* ::Input:: *)
(*pointC = Point[coordC]*)


(* ::Input:: *)
(*Graphics[{lineAB, lineAC,lineBC, circleACE, circleBCD, Cyan, Triangle[{coordA, coordB, coordC}], PointSize[Large], Red, pointA, pointB, Green, pointC}]*)
