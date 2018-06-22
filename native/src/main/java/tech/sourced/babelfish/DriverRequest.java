package tech.sourced.babelfish;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.io.IOException;

/**
 * Class for the C/C++ driver request.
 */


class DriverRequest {

    static class RequestLoadException extends IOException {
        RequestLoadException(Throwable e) {
            super(e);
        }
    }
    public String action;
    public String language;
    public String languageVersion;
    public String content;
    @JsonIgnoreProperties(ignoreUnknown = true)

    public DriverRequest() {
        // Dummy constructor, jackson needs this
    }
    // TODO: check if Jackson needs this
    DriverRequest(String action, String language, String languageVersion, String content) {
        this.content = content;
        this.action = action;
        this.language = language;
        this.languageVersion = languageVersion;
    }

    static DriverRequest load(String in) throws RequestLoadException {
        // NOTE: If we add new protocols this need to be declouped from jackson through an
        // intermediate interface (IExchangeFormatReader) like DriverResponse is, but for now with a single protocol is overkill
        // to add more layers

        // TODO: check if this "new" can be avoided since it will be called in the loop
        ObjectMapper mapper = new ObjectMapper();

        try {
            return mapper.readValue(in, DriverRequest.class);
        } catch (IOException e) {
            throw new DriverRequest.RequestLoadException(e);
        }
    }
}
