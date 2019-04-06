package edgar

var tableTemplate = `
<html>
<body>
  <table border="1">
    <tr>
      <th>
        "Filed"
      </th>
      <th>
        "Equity"
      </th>
      <th>
        "Shares"
      </th>
      <th>
        "Revenue"
      </th>
      <th>
        "CoRev"
      </th>
      <th>
        "OpsExp"
      </th>
      <th>
        "OpsInc"
      </th>
      <th>
        "NetInc"
      </th>
      <th>
        "LDebt"
      </th>
      <th>
        "SDebt"
      </th>
      <th>
        "Cash"
      </th>
    </tr>
    {{ range $date, $report := . }}
      <tr>
        <th>
          {{ $date }}
        </th>
        <th>
          {{ printf "%.2f" $report.FinData.Bs.Equity }}
        </th>
        <th>
          {{ printf "%.2f" $report.FinData.Entity.ShareCount }}
        </th>
        <th>
          {{ printf "%.2f" $report.FinData.Ops.Revenue }}
        </th>
        <th>
          {{ printf "%.2f" $report.FinData.Ops.CostOfSales }}
        </th>
        <th>
          {{ printf "%.2f" $report.FinData.Ops.OpExpense }}
        </th>
        <th>
          {{ printf "%.2f" $report.FinData.Ops.OpIncome }}
        </th>
        <th>
          {{ printf "%.2f" $report.FinData.Ops.NetIncome }}
        </th>
        <th>
          {{ printf "%.2f" $report.FinData.Bs.LDebt }}
        </th>
        <th>
          {{ printf "%.2f" $report.FinData.Bs.SDebt }}
        </th>
        <th>
          {{ printf "%.2f" $report.FinData.Bs.Cash }}
        </th>
      </tr>
   {{ end }}
  </table>
</body>
</html>
`
