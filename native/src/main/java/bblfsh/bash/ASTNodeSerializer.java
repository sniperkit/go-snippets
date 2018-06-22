package bblfsh.bash;

import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.databind.SerializerProvider;
import com.fasterxml.jackson.databind.ser.std.StdSerializer;
import com.intellij.lang.ASTNode;
import com.intellij.psi.tree.TokenSet;

import java.io.IOException;
import java.util.List;

/**
 * Custom Jackson serializer for com.intellij.lang.ASTNode;
 */
public class ASTNodeSerializer extends StdSerializer<ASTNode> {

    public ASTNodeSerializer() {
        this(null);
    }

    public ASTNodeSerializer(Class<ASTNode> t) {
        super(t);
    }

    @Override
    public void serialize(ASTNode root, JsonGenerator jG, SerializerProvider provider) throws IOException {
        jG.writeStartObject();

        final String type = root.getElementType().toString();
        jG.writeStringField("elementType", type);

        final int start = root.getStartOffset();
        jG.writeNumberField("startOffset", start);

        final int length = root.getTextLength();
        jG.writeNumberField("textLength", length);

        final String text = root.getText();
        jG.writeStringField("text", text);

        jG.writeFieldName("children");
        jG.writeStartArray();
        final TokenSet filter = null;
        for (ASTNode child : root.getChildren(filter)) {
            serialize(child, jG, provider);
        }
        jG.writeEndArray();

        jG.writeEndObject();
    }

}
