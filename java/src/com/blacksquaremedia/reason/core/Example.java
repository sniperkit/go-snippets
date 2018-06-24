package com.blacksquaremedia.reason.core;

// An example instance.
public interface Example {
  // Returns the value of a given featureName. Must be either a String or an Integer for
  // categorical features or a Double/Interger for numerical features.
  public Object getExampleValue(java.lang.String featureName);
}
