package bblfsh.bash;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.module.SimpleModule;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import com.intellij.lang.ASTNode;

public class Registry {

    private static ObjectMapper mapper = null;

    private Registry() {}

    public static ObjectMapper objectMapper() {
        if (mapper == null) {
            initObjectMapper();
        }
        return mapper;
    }

    private static void initObjectMapper() {
        mapper = new ObjectMapper();

        mapper.disable(SerializationFeature.FAIL_ON_EMPTY_BEANS);
        mapper.disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES);
        mapper.enable(DeserializationFeature.ACCEPT_EMPTY_STRING_AS_NULL_OBJECT);
        mapper.setSerializationInclusion(JsonInclude.Include.NON_NULL);

        final SimpleModule module = new SimpleModule();
        module.addSerializer(ASTNode.class, new ASTNodeSerializer());
        mapper.registerModule(module);
    }

}
