/* global window, __BUILD__ */


// mock analytics package when not in production
if (!__BUILD__.PRODUCTION || __BUILD__.SEGMENT_KEY == "") {
  const noop = () => {};
  window.analytics = {
    track: noop,
    identify: noop,
    page: noop,
    group: noop,
    alias: noop,
    ready: noop,
  };
}

const analytics = window.analytics;
const appProps = JSON.parse(__BUILD__.SEGMENT_APP_PROPERTIES);

// analytics shim to inject global values
// further reading: https://segment.com/docs/sources/website/analytics.js
export default {
  track(event, properties, options, callback) {
    return analytics.track(event, Object.assign(appProps, properties), options, callback);
  },
  identify(userId, traits, options, callback) {
    return analytics.identify(userId, traits, options, callback);
  },
  page(category, name, properties, options, callback) {
    return analytics.page(category, name, Object.assign(appProps, properties), options, callback);
  },
  group(groupId, traits, options, callback) {
    return analytics.group(groupId, Object.assign(appProps, traits), options, callback);
  },
  alias(userId, previousId, options, callback) {
    return analytics.alias(userId, previousId, options, callback);
  },
  ready(callback) {
    return analytics.ready(callback);
  },
};
