<html>
  <h1>WattView</h1>
<h2>Summary</h2>
WattView - a collection of Go language command line routines written on the Raspberry Pi for the WattTime (watttime.org) API to read real-time and historic electric power grid information on the USA (and eventually European) Grids.<br /><br />
Written by Maurice Bizzarri, Bizzarri Software, January 2019<br />
See https://watttime.org for more information on this system<br />
  <h2>Getting Started</h2>
  Just clone the repository and compile the go routines with go install.  There is also a script called INSTALL that can be run after cloning <br />
  <h2>Prerequisites</h2>
  A working GO environment.  No fancy features used so older GO installs should work.<br />
  <h2>Documentation on the routines</h2>
All routines support the -h flag to document the parameters<br />
<p>
emissions - get the emissions status for a Balancing Authority<br />
griddata - get the detailed grid data for a specific BA over a specific time period.  This may be restricted according to the type of account you have.<br />
gridregion - get the grid region (BA) for a specific longitude/latitude pair<br />
makeacct  - make an account for WattTime<br />
</p>
<p>- makeacct will make a free account on WattTime.org.  It will also write the account name, password, and other info you supplied to WattTime.org into a file ($HOME/.WattTime/account) that the other programs will look for and use.  You can also create this file from scratch if you already have an account.  The file is documented later in this README.</p>
<p>- gridregion takes a latitude/longitude pair and returns the short abbreviation for the balancing authority.  It also writes the abbreviation into the location file in .WattTime directory.  This will also be used by other programs.</p>
<p>- emissions will return the emissions status for a specific Balancing Authority.
You can specify a BA on the command line or use the one already in the location file.  There is a quiet mode that will allow use of this function in scripts.  It returns 0 if the grid is green, else it returns 1.</p>
<p>- griddata returns detailed grid power supply history information on 5 minute increments.  In a free account this is limited to CAISO_ZP26 BA.  It uses the account file but can be overridden on the command line. You can also specify a file name to write the grid data to instead. It can write a JSON formatted file or a CSV (comma separated value) file.</p>
<p>- passrecover allows password recovery of an existing account.  WattTime will send a recovery link to the email address associated with the account.</p>
<p>
<h2>Author</h2>
Maurice Bizzarri<br />
Bizzarri Software<br />
maurice@bizzarrisoftware.com<br />
https://bizzarrisoftware.com<br />
</p>
<h2>License</h2>
This project is licensed under the MIT License - see the LICENSE.md file for details<br />
<h2>Files Created</h2>
$HOME/.WattTime/account<br />
Contains information on account - written by makeacct<br />
Example File: <br />
{"username":"account-name","password":"account-password","email":"email address","org":"organization name"}<br /><br />
$HOME/.WattTime/ba<br />
Contains location searched for in gridregion.  Will be used by griddata and emissions. Can be overridden by using command line arguments.<br />
Example File:<br />
CAISO_ZP26<br />
</html>
