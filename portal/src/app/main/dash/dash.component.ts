import {Component, OnInit} from '@angular/core';
import {ApiService} from '../../api.service';
import {TabRef} from "../tabs/tabs.component";

@Component({
  selector: 'app-dash',
  templateUrl: './dash.component.html',
  styleUrls: ['./dash.component.scss']
})
export class DashComponent implements OnInit {

  constructor(private as: ApiService, private tab: TabRef) {
    tab.name = '仪表盘';
  }

  ngOnInit(): void {

  }

}
