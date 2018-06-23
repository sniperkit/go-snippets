package tech.sourced.babelfish;

import com.fasterxml.jackson.core.JsonEncoding;
import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.core.util.DefaultPrettyPrinter;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import com.fasterxml.jackson.databind.module.SimpleModule;
import org.eclipse.cdt.core.dom.ast.IASTTranslationUnit;

import java.io.ByteArrayOutputStream;
import java.io.IOException;

class TranslationUnitJSONMapper implements IExchangeFormatWritter {

    final JsonGenerator generator;
    final JsonFactory jsonFactory = new JsonFactory();
    final ObjectMapper mapper = new ObjectMapper();
    private ByteArrayOutputStream byteOutputStream;

    TranslationUnitJSONMapper(boolean prettyPrint, ByteArrayOutputStream byteOutput) throws IOException {
        this.byteOutputStream = byteOutput;

        generator = jsonFactory.createGenerator(byteOutputStream, JsonEncoding.UTF8);
        if (prettyPrint) {
            generator.setPrettyPrinter(new DefaultPrettyPrinter());
            mapper.enable(SerializationFeature.INDENT_OUTPUT);
        }
        mapper.disable(SerializationFeature.FAIL_ON_EMPTY_BEANS);
        mapper.disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES);
        mapper.enable(DeserializationFeature.ACCEPT_EMPTY_STRING_AS_NULL_OBJECT);
        SimpleModule module = new SimpleModule();
        module.addSerializer(TranslationUnit.class, new TranslationUnitSerializer());
        mapper.registerModule(module);
    }

    public void writeValue(DriverResponse response) throws IOException {
        mapper.writeValue(generator, response);
    }

    public ByteArrayOutputStream getByteOutputStream() {
        return byteOutputStream;
    }
}
