package bblfsh.bash;

import com.fasterxml.jackson.databind.ObjectMapper;

public class Main {

    public static void main(String args[]) {
        final ObjectMapper mapper = Registry.objectMapper();
        final RequestReader reader = new RequestReader(System.in, mapper);
        final ResponseWriter writer = new ResponseWriter(System.out, mapper);
        final Driver driver = new Driver(reader, writer);

        try {
            driver.run();
        } catch (CloseException e) {
            System.exit(0);
        } catch (DriverException e) {
            System.exit(1);
        }
    }

}
