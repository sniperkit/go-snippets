package bblfsh.bash;

/**
 * Class for the java driver request.
 */
public class Request {
    public String content;

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;

        Request request = (Request) o;

        return content != null ? content.equals(request.content) : request.content == null;
    }

    @Override
    public int hashCode() {
        return content != null ? content.hashCode() : 0;
    }
}
