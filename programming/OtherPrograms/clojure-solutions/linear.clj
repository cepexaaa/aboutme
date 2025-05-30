(defn hop-hey-la-la-ley
  [op]
  (fn [v1 v2]
    (mapv op v1 v2)))

(def v+ (hop-hey-la-la-ley +))
(def v- (hop-hey-la-la-ley -))
(def v* (hop-hey-la-la-ley *))
(def vd (hop-hey-la-la-ley /))


(defn scalar [v1 v2] (apply + (mapv * v1 v2)))
(defn vect [v1 v2]
  (let [[x1 y1 z1] v1 [x2 y2 z2] v2]
    [(- (* y1 z2) (* z1 y2))
     (- (* z1 x2) (* x1 z2))
     (- (* x1 y2) (* y1 x2))]))
(defn v*s [v s] (mapv #(* % s) v))

(defn dav_matrix [op] (fn [m1 m2] (mapv (fn [v1 v2] (op v1 v2)) m1 m2)))
(def m+ (hop-hey-la-la-ley v+))
(def m- (dav_matrix v-))
(def m* (dav_matrix v*))
(def md (dav_matrix vd))

(defn transpose [m] (apply mapv vector m))
(defn m*s [m s] (mapv #(mapv (fn [p1] (* p1 s)) %) m))
(defn m*v [m v] (mapv #(scalar % v) m))
(defn m*m [m1 m2] (mapv (fn [row] (m*v (transpose m2) row)) m1))

(defn dav [op]
  (fn [m1 m2]
    (mapv (fn [m1-slice m2-slice]
            (mapv (fn [m1-matrix m2-matrix]
                    (op m1-matrix m2-matrix))
                  m1-slice m2-slice))
          m1 m2)))
(def c4+ (dav m+))
(def c4- (dav m-))
(def c4* (dav m*))
(def c4d (dav md))

;(defn parseObjectPostfix [expr-str]
;  ;(println expr-str)
;  (letfn [(parse [expr]
;            (cond
;              (list? expr) (apply (operation-object-map (last expr)) (map parse (pop expr)))
;              (number? expr) (Constant expr)
;              :else (Variable (str expr))))]
;    (parse (read-string expr-str))))


;(toStringPostfix (parseObjectPostfix (2 negate)))