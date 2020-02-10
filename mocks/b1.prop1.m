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


(* ::Input:: *)
(*Manipulate[*)
(*r =EuclideanDistance[a,b];*)
(*ab = Line[{a,b}];*)
(*ace = Circle[a, r];*)
(*bde = Circle[b, r];*)
(*pts =  {x, y} /. Quiet[Solve[{x,y} \[Element] ace && {x,y} \[Element] bde, {x,y}]];*)
(*c = Last[pts];*)
(*Graphics[*)
(*{PointSize[0.03], Red, Point[a], Point[b],*)
(*Black, ab, ace,bde,*)
(*Green, Point[c],*)
(*Cyan, Triangle[{a, b, c}],*)
(*Gray, Text[Style["A", 20], a, {2,0}],*)
(*Text[Style["B", 20], b, {-2,0}],*)
(*Text[Style["C", 20], c, {0,-1.5}\.00]*)
(*},*)
(*PlotRange -> 4 (*,Axes\[Rule]True*)],*)
(*{{a,{-1,0}},Locator},*)
(*{{b, {1, 0}}, Locator}]*)
