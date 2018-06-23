package bblfsh.bash;

import java.util.ArrayList;
import com.intellij.lang.ASTNode;

/**
 * Class for the java driver response
 */
public class Response {
    public String status;
    public ArrayList<String> errors = new ArrayList<String>(0);
    public ASTNode ast;

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;

        Response response = (Response) o;

        if (status != null ? !status.equals(response.status) : response.status != null)
            return false;
        if (errors != null ? !errors.equals(response.errors) : response.errors != null)
            return false;
        return ast != null ? ast.equals(response.ast) : response.ast == null;
    }

    @Override
    public int hashCode() {
        int result = status != null ? status.hashCode() : 0;
        result = 31 * result + (errors != null ? errors.hashCode() : 0);
        result = 31 * result + (ast != null ? ast.hashCode() : 0);
        return result;
    }
}
