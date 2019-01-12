<html>
  <h1>WattTime</h1>
<h2>Summary</h2>
These Go stand-alone programs create a simple command line interface for the WattTime system to read real-time and historic electric power grid information on the USA (and eventually European) Grids.<br />
Written by Maurice Bizzarri, Bizzarri Software, January 2019<br />
see https://watttime.org for more information on this system<br />
  <h2>Getting Started</h2>
  Just clone the files and compile the go routines with go install <br />
  <h2>Prerequisites</h2>
  A working GO environment.  No fancy features used so older GO installs should work.<br />
  
<p>
emissions - get the emissions status for a Balancing Authority<br />
griddata - get the detailed grid data for a specific BA over a specific time period.  This may be restricted according to the type of account you have.<br />
gridregion - get the grid region (BA) for a specific longitude/latitude pair<br />
makeacct  - make an account for WattTime<br />
</p>
<p>- makeacct will make a free account on WattTime.org.  It will also write the account name, password, and other info you supplied to WattTime.org into a file ($HOME/.WattTime/account) that the other programs will look for and use.  You can also create this file from scratch if you already have an account.  The file is documented later in this README.</p>
<p>- gridregion takes a latitud/longitude pair and returns the short abbreviation for the balancing authority.  It also writes the abbreviation into the location file in .WattTime directory.  This will also be used by other programs.</p>
<p>- emissions will return the emissions status for a specific Balancing Authority.
You can specify a BA on the command line or use the one already in the location file.</p>
<p>- griddata returns detailed grid power supply history information on 5 minute increments.  In a free account this is limited to CAISO_ZP26 BA.  It uses the account file but can be overridden on the command line.</p>
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
{"username":"account-name","password":"account-password","email":"email address","org":"organization name"}<br />
$HOME/.WattTime/location<br />
Contains location searched for in gridregion<br />
Example File<br />
CAISO_ZP26<br />
</html>
