(ns babelfish-clojure-driver.parse-test
  (:require [clojure.test :refer :all]
            [babelfish-clojure-driver.parse :refer :all]
            [babelfish-clojure-driver.core :refer [form-ast]]))

(defn- assert-comment
  [node text]
  (is (= (-> node :exprs first) text)))

(deftest comment-node-test
  (testing "returns a valid comment node"
    (let [node (comment-node "foo" 1 2 6)]
      (is (= (:op node) :comment))
      (is (nil? (:form node)))
      (is (= (:meta node) {:line 1
                           :col 2
                           :end-line 1
                           :end-col 6}))
      (is (= (-> node :exprs count) 1))
      (assert-comment node "foo")
      (is (= (:children node) [:exprs])))))

(def test-src
  ";; this is a comment

  (def a b) ;; another one

  (def \"this \\\"is ; not a comment
   fooo bar baz\")
  
  (def c \\;)
  
  ;; yet another")

(deftest extract-comments-test
  (testing "extracts all comments"
    (let [comments (extract-comments test-src)]
      (is (= (count comments) 3))
      (let [[c1 c2 c3] comments]
        (assert-comment c1 "this is a comment")
        (assert-comment c2 "another one")
        (assert-comment c3 "yet another")))))

(deftest parse-defn-test
  (testing "defn should be expanded"
    (let [form '(defn foo "documentation" [a b] a)
          ast (form-ast form)]
      (is (= :def (-> ast :op)))
      (is (= "documentation" (-> ast :meta :doc)))
      (let [init (:init ast)]
        (is (= :fn (-> init :op)))
        (is (= 2 (-> init :max-fixed-arity)))
        (is (= 1 (-> init :methods count))))))

  (testing "defn should expand multi-signature"
    (let [form '(defn foo "documentation" ([a b] a) ([a b c] a))
          ast (form-ast form)
          init (:init ast)]
      (is (= :fn (-> init :op)))
      (is (= 3 (-> init :max-fixed-arity)))
      (is (= 2 (-> init :methods count)))))

  (testing "defn should accept attr maps"
    (let [form '(defn foo "documentation" {:b 2} ([a b] a) ([a b c] a) {:a 1})
          ast (form-ast form)
          init (:init ast)]
      (is (= :fn (-> init :op)))
      (is (= 3 (-> init :max-fixed-arity)))
      (is (= 1 (-> ast :meta :a)))
      (is (= 2 (-> ast :meta :b)))
      (is (= 2 (-> init :methods count))))))
