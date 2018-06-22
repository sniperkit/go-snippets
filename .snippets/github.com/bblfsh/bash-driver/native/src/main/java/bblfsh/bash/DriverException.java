package bblfsh.bash;

public class DriverException extends Exception {

    public DriverException(final String msg) {
        super(msg);
    }

    public DriverException(final String msg, final Throwable cause) {
        super(msg, cause);
    }

}
