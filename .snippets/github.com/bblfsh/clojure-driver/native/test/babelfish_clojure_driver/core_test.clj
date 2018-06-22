(ns babelfish-clojure-driver.core-test
  (:require [clojure.test :refer :all]
            [clojure.data.json :as json]
            [babelfish-clojure-driver.core :refer :all]))

(defn- err?
  [result status]
  (do (is (= (:status result) status))
      (is (nil? (:ast result)))
      (is (= (count (:errors result)) 1))))

(defn- ok?
  [result status]
  (do (is (= (:status result) status))
      (is (not (nil? (:ast result))))
      (is (= (count (:errors result)) 0))))

(defn- read-str
  [s]
  (json/read-str s :key-fn keyword :eof-error? false :eof-value :end))

(defn- stream
  [req]
  (java.io.BufferedReader. (java.io.StringReader. req)))

(deftest test-parse
  (testing "with valid source"
    (let [result (parse "(defn foo [a b] (+ a b))")
          m (-> result :ast :meta)]
      (ok? result :ok)))

  (testing "includes positional info"
    (let [result (parse "(def a\n1)")
          m (-> result :ast first :meta)]
      (is (= (:status result) :ok))
      (is (= (:line m) 1))
      (is (= (:column m) 1))
      (is (= (:end-line m) 2))
      (is (= (:end-column m) 3))))

  (testing "with invalid source"
    (let [result (parse "(println (+ a 1)")
          m (-> result :ast :meta)]
      (err? result :error))))

(deftest test-process-req
  (testing "with non json input"
    (let [result (process-req (stream "fooooo"))
          output (read-str result)]
      (err? output "fatal")))

  (testing "with valid input"
    (let [req (stream (json/write-str {:content "(defn foo [a b]
                                                  (+ a b))"}))
          result (process-req req)
          output (read-str result)]
      (ok? output "ok"))))

(deftest test-clean-ast
  (testing "cleaning ast"
    (let [ast [{:env true
                :raw-forms true
                :other {:env true
                        :raw-forms true
                        :var true
                        :other true}
                :var true}
               {[(symbol "quote") (symbol "foo")] {:env true
                                                   :raw-forms true
                                                   :var true
                                                   :other true}}]
          clean (clean-ast ast)]
      (is (= [{:other {:other true}}
              {"'foo" {:other true}}]
             clean)))))
