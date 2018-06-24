package com.blacksquaremedia.reason.classification;

import java.math.BigInteger;
import com.blacksquaremedia.reason.CoreProtos;
import com.blacksquaremedia.reason.UtilProtos;
import com.blacksquaremedia.reason.core.Example;
import com.google.protobuf.ProtocolStringList;
import net.jpountz.xxhash.XXHash64;
import net.jpountz.xxhash.XXHashFactory;

public class Hoeffding {
  private HoeffdingProtos.Tree tree;
  private XXHashFactory hasher;

  // The standard constructor wraps a Tree object
  public Hoeffding(HoeffdingProtos.Tree tree) {
    this.tree = tree;
    this.hasher = XXHashFactory.fastestInstance();
  }

  // Lookup traverses the tree and returns StreamStats
  // that best represent the passed example.
  public UtilProtos.Vector lookup(Example example) {
    if (!this.tree.hasModel()) {
      return UtilProtos.Vector.newBuilder().build();
    }

    // Get the model and the root node.
    CoreProtos.Model model = this.tree.getModel();

    // Traverse the tree starting at the root node.
    HoeffdingProtos.Node node = this.traverse(model, example, this.tree.getRoot());
    if (node == null) {
      return UtilProtos.Vector.newBuilder().build();
    }
    return node.getStats();
  }

  private HoeffdingProtos.Node traverse(CoreProtos.Model model, Example example, long nodeRef) {
    HoeffdingProtos.Node node = this.tree.getNodes((int) (nodeRef - 1));

    // If the node is a split node ....
    if (node.hasSplit()) {
      HoeffdingProtos.SplitNode split = node.getSplit();

      // Lookup the associated feature by name.
      CoreProtos.Feature feature = model.getFeaturesOrThrow(split.getFeature());

      // Get the example value for the feature, stop traversing if
      // we cannot descend further.
      Object obj = example.getExampleValue(feature.getName());
      if (obj == null) {
        return node;
      }

      // Lookup the childRef based in the example value and the feature type
      // (either categorical or numerical).
      long childRef = 0;
      switch (feature.getKind()) {
        case CATEGORICAL:
          childRef = this.findChildRefCategorical(split.getChildren(), feature, obj);
          break;
        case NUMERICAL:
          childRef = this.findChildRefNumerical(split.getChildren(), obj, split.getPivot());
          break;
      }

      // If a child can be found, continue traversal.
      if (childRef > 0) {
        return this.traverse(model, example, childRef);
      }
    }

    // We must be at a leaf node or we cannot traverse any further.
    return node;
  }

  // Looks up a child ref using a categorical feature's value.
  private long findChildRefCategorical(
      HoeffdingProtos.SplitNode.Children nodes, CoreProtos.Feature feature, Object obj) {
    switch (feature.getStrategy()) {
      case IDENTITY:
        // If the feature stategy is IDENTITY, try to cast/convert the value into a number
        // and use the result as the index to lookup child node.
        if (obj instanceof Number) {
          return this.findChildRef(nodes, ((Number) obj).intValue());
        } else if (obj instanceof String) {
          return this.findChildRef(nodes, Integer.parseInt((String) obj));
        }
        break;
      case VOCABULARY:
      case EXPANDABLE:
        // If the feature stategy is VOCABULARY/EXPANDABLE, obtain the string
        // value first.
        String value = obj.toString();
        ProtocolStringList vocabulary = feature.getVocabularyList();

        // Check if that string is included in the vocabulary; use the index if so.
        int index = vocabulary.indexOf(value);
        if (index > -1) {
          return this.findChildRef(nodes, index);
        }

        // If the string isn't stored in the vocabulary, check if this feature
        // supports bucket hashing.
        int numBuckets = feature.getHashBuckets();
        if (numBuckets > 0) {

          // Calculate the index as HASH(value) % NUM_BUCKETS + SIZE(VOCABULARY)
          byte[] data = value.getBytes();
          XXHash64 hash = this.hasher.hash64();
          long hashed = hash.hash(data, 0, data.length, 0);
          BigInteger bigInt = new BigInteger(Long.toUnsignedString(hashed));
          BigInteger bucket = bigInt.remainder(BigInteger.valueOf(numBuckets));

          return this.findChildRef(nodes, bucket.intValue() + vocabulary.size());
        }
        break;
    }
    return 0;
  }

  // Looks up a child node using a numerical feature's value.
  private long findChildRefNumerical(
      HoeffdingProtos.SplitNode.Children nodes, Object obj, double pivot) {
    if (obj instanceof Number) {
      double value = ((Number) obj).doubleValue();
      if (value < pivot) {
        return this.findChildRef(nodes, 0);
      } else {
        return this.findChildRef(nodes, 1);
      }
    }
    return 0;
  }

  private long findChildRef(HoeffdingProtos.SplitNode.Children children, int index) {
    if (index > -1) {
      if (index < children.getDenseCount()) {
        return children.getDense(index);
      } else if (children.getSparseCount() != 0) {
        return children.getSparseOrDefault((long) (index), 0);
      }
    }
    return 0;
  }
}
