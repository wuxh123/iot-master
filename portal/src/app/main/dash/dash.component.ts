import { Component, OnInit } from '@angular/core';
import {CaService} from '../../ca.service';

@Component({
  selector: 'app-dash',
  templateUrl: './dash.component.html',
  styleUrls: ['./dash.component.scss']
})
export class DashComponent implements OnInit {

  constructor(private ca: CaService) {
  }

  ngOnInit(): void {

  }

}
