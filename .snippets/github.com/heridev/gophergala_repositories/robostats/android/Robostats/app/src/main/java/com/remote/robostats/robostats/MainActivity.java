package com.remote.robostats.robostats;

import android.content.Intent;
import android.os.AsyncTask;
import android.os.Bundle;
import android.support.v7.app.ActionBarActivity;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.Button;

import org.xmlpull.v1.XmlPullParser;
import org.xmlpull.v1.XmlPullParserException;
import org.xmlpull.v1.XmlPullParserFactory;

import java.io.BufferedInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.net.HttpURLConnection;
import java.net.URL;



public class MainActivity extends ActionBarActivity {

    public final static String robostatsUserName = "user";
    public final static String robostatsnPassword = "pass";
    public final static String apiURL = "?";

    private class emailVerificationResult {
        public String status;
        public String result;
    }

    private class CallAPI extends AsyncTask<String, String, String> {


        @Override
        protected String doInBackground(String... params) {
            String urlString=params[0]; // URL to call
            String resultToDisplay = "";
            InputStream in = null;
            emailVerificationResult result = null;

            // HTTP Get
            try {
                URL url = new URL(urlString);
                HttpURLConnection urlConnection = (HttpURLConnection) url.openConnection();
                in = new BufferedInputStream(urlConnection.getInputStream());
            } catch (Exception e ) {
                System.out.println(e.getMessage());
                return e.getMessage();
            }

            // Parse XML
            XmlPullParserFactory pullParserFactory;
            try {
                pullParserFactory = XmlPullParserFactory.newInstance();
                XmlPullParser parser = pullParserFactory.newPullParser();
                parser.setFeature(XmlPullParser.FEATURE_PROCESS_NAMESPACES, false);
                parser.setInput(in, null);
                result = parseXML(parser);
            } catch (XmlPullParserException e) {
                e.printStackTrace();
            } catch (IOException e) {
                e.printStackTrace();
            }

            // Logic to determine if the login is invalid, or valid
            if (result != null ) {
                if( Integer.parseInt(result.status) >= 300) {
                    resultToDisplay = "Invalid login credentials";
                }
                else {
                    resultToDisplay = "";
                }

            }
            else {
                resultToDisplay = "Exception Occured";
            }

            return resultToDisplay;

        }

        protected void onPostExecute(String result) {
            Button login_button = (Button) findViewById(R.id.login_button);
            login_button.setText(R.string.login_button);

            login_button.setOnClickListener(new View.OnClickListener() {
                public void onClick(View v) {
                    Intent i = new Intent(MainActivity.this, com.remote.robostats.robostatsremote.UserMenuScreen.class);
                    MainActivity.this.startActivity(i);
                }
            });
        }

        private emailVerificationResult parseXML( XmlPullParser parser ) throws XmlPullParserException, IOException {
            int eventType = parser.getEventType();
            emailVerificationResult result = new emailVerificationResult();

            while( eventType!= XmlPullParser.END_DOCUMENT) {
                String name = null;
                switch(eventType)
                {
                    case XmlPullParser.START_TAG:
                        name = parser.getName();
                        if( name.equals("Error")) {
                            System.out.println("");
                        }
                        else if ( name.equals("")) {
                            result.status = parser.nextText();
                        }
                        else if (name.equals("")) {
                            result.result = parser.nextText();
                        }
                        break;
                    case XmlPullParser.END_TAG:
                        break;
                }

                eventType = parser.next();
            }
            return result;
        }


    } // end CallAPI


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);


    }


    @Override
    public boolean onCreateOptionsMenu(Menu menu) {

        getMenuInflater().inflate(R.menu.menu_main, menu);
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {

        int id = item.getItemId();

        if (id == R.id.action_settings) {
            return true;
        }

        return super.onOptionsItemSelected(item);
    }
}
