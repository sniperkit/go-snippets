try {
    throw new RangeError();
} catch (e) {
    console.log("catched");
} finally {
    console.log("always");
}
