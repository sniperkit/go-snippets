package bblfsh.bash;

import java.io.IOException;

// CloseException signals that stream is cleanly closed.
public class CloseException extends IOException {

    public CloseException(final String msg) {
        super(msg);
    }

    public CloseException(final String msg, final Throwable cause) {
        super(msg, cause);
    }

}
