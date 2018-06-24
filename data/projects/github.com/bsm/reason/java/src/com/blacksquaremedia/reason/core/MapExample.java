package com.blacksquaremedia.reason.core;

import java.util.Map;

public class MapExample implements Example {
  private java.util.Map<String,Object> map;

  MapExample(java.util.Map<String,Object> map)    {
    this.map = map;
  }

  public Object getExampleValue(java.lang.String featureName) {
    return this.map.get(featureName);
  }
}
