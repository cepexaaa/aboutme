;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;Functional expression;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(def constant constantly)
(defn variable [name] (fn [vars] (get vars name)))
(defn operation [f] (fn [& args] (fn [vars] (apply f (mapv #(% vars) args)))))
(def add (operation +))
(def subtract (operation -))
(def multiply (operation *))
(def divide (operation (fn
                         ([a] (/ (double a)))
                         ([a b] (/ (double a) (double b))))))
(def negate (operation -))
(def sinh (operation #(Math/sinh %)))
(def cosh (operation #(Math/cosh %)))
(def operation-map {'+ add, '- subtract, '* multiply, '/ divide, 'negate negate, 'sinh sinh, 'cosh cosh})
(defn parser [expr-str, constant-f-o, variable-f-o, op-map-f-o]
  (letfn [(parse [expr]
            (cond
              (list? expr) (apply (op-map-f-o (first expr)) (map parse (rest expr)))
              (number? expr) (constant-f-o expr)
              :else (variable-f-o (str expr))))]
    (parse (read-string expr-str))))

(defn parseFunction [str-exp] (parser str-exp constant variable operation-map))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;Object expression;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(defn proto-get
  ([obj key] (proto-get obj key nil))
  ([obj key default]
   (or (get obj key)
       (when-let [prototype (:prototype obj)]
         (proto-get prototype key default))
       default)))

(defn proto-call
  [this key & args]
  (apply (proto-get this key) this args))

(defn field [key]
  (fn
    ([this] (proto-get this key))
    ([this def] (proto-get this key def))))

(defn method [key]
  (fn [this & args] (apply proto-call this key args)))

(defn constructor [ctor prototype]
  (fn [& args] (apply ctor (assoc prototype :prototype prototype) args)))

(def evaluate (method :evaluate))
(def toString (method :toString))
(def toStringPostfix (method :toStringPostfix))

(defn protoInterface [eval toStr toStrPost]
  {:evaluate eval
   :toString toStr
   :toStringPostfix toStrPost})

(def Constant (let [value (field :val)]
                (constructor
                  (fn [this val]
                    (assoc this :val val))
                  (protoInterface
                    (fn [this _] (value this))
                    (fn [this] (str (value this)))
                    (fn [this] (str (value this)))))))

(def Variable (let [nameVar (field :name)]
                (constructor
                  (fn [this name]
                    (assoc this :name name))
                  (protoInterface
                    (fn [this vars] (get vars (str (clojure.string/lower-case (subs (nameVar this) 0 1)))))
                    (fn [this] (nameVar this))
                    (fn [this] (nameVar this))))))


(defn Operation [operation-sign operator]
  (let [operation-name (field :operationSign)
        arguments (field :args)]
    (fn [& args]
      {:args args
       :prototype
       {:operationSign operation-sign
        :evaluate (fn [this vars] (apply operator (map #(evaluate % vars) (arguments this))))
        :toString (fn [this] (format "(%s %s)" (operation-name this) (clojure.string/join " " (map toString (arguments this)))))
        :toStringPostfix (fn [this] (format "(%s %s)" (clojure.string/join " " (map toStringPostfix (arguments this))) (operation-name this)))}})))

(def Add (Operation "+" +))
(def Subtract (Operation "-" -))
(def Multiply (Operation "*" *))
(def Divide (Operation "/" (fn ([a] (/ (double a))) ([a b] (/ (double a) (double b))))))
(def Negate (Operation "negate" -))
(def Pow (Operation "pow" #(Math/pow % %2)))
(def Log (Operation "log" (fn [a b] (/ (Math/log (Math/abs b)) (Math/log (Math/abs a))))))
(def Max (Operation "max" max))
(def Min (Operation "min" min))
(def operation-object-map {'+ Add, '- Subtract, '* Multiply, '/ Divide, 'negate Negate, 'pow Pow, 'log Log, 'max Max, 'min Min})
(defn parseObject [str-exp] (parser str-exp Constant Variable operation-object-map))


;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;Combination parser;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
(load-file "parser.clj")

(def constant_comb (+map #(Constant (read-string %))
                     (+str (+seqf (fn [s t] (if (#{\- \+} s) (cons s t) t))
                                  (+opt (+char "-"))
                                  (+plus *digit)))))

(def variable_comb (+map #(Variable %) (+string "xyzXYZ")))

(def *operator
  (let [ops (map name (keys operation-object-map))]
    (+map #(do
             (println %)
             (operation-object-map (symbol %)))
          (+seqf cons (+char (first ops))
                 (apply +seq (rest ops))))))

(def *op
  (+seqf
    (fn [operands operator] (apply operator operands))
    (+ignore (+char "(")) *ws
    (+plus (+seqn 0 (+or constant_comb variable_comb (delay *op)) *ws))
    *ws *operator *ws (+ignore (+char ")"))))


(def parseObjectPostfix (_parser (+seqn 0 *ws (+or constant_comb variable_comb *op) *ws)))

