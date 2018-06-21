use std::collections::LinkedList;
use std::collections::BTreeMap;

fn main() {
    fn linklist() {
        let mut list1 = LinkedList::new();
        list1.push_back('a');
        let mut list2 = LinkedList::new();
        list2.push_back('b');
        list2.push_back('c');

        list1.append(&mut list2);

        let mut iter = list1.iter();
        assert_eq!(iter.next(), Some(&'a'));
        assert_eq!(iter.next(), Some(&'b'));
        assert_eq!(iter.next(), Some(&'c'));
        assert!(iter.next().is_none());

        assert!(list2.is_empty());
    }

    fn btreemap() {
        let mut map: BTreeMap<&str, usize> = BTreeMap::new();
        map.entry("liwei").or_insert(100);

        assert_eq!(map["liwei"], 101);
    }

    linklist();
    btreemap();
}
