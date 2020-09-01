import {Component, OnInit, ViewContainerRef} from '@angular/core';
import {ApiService} from "../../api.service";
import {NzDrawerService, NzModalService} from "ng-zorro-antd";
import {LinkEditComponent} from "../link-edit/link-edit.component";
import {LinkDetailComponent} from "../link-detail/link-detail.component";

@Component({
  selector: 'app-link',
  templateUrl: './link.component.html',
  styleUrls: ['./link.component.scss']
})
export class LinkComponent implements OnInit {
  links: [];

  constructor(private as: ApiService, private drawer: NzDrawerService) {
  }

  ngOnInit(): void {
    this.load();
  }

  load(): void {
    this.as.get('links').subscribe(res => {
      if (res.ok) {
        this.links = res.data;
      }
    });
  }

  edit(l?): void {
    this.drawer.create({
      nzTitle: l ? '编辑' : '创建',
      nzMaskClosable: false,
      nzWidth: 500,
      nzContent: LinkEditComponent,
      nzContentParams: {
        link: l || {}
      }
    });
  }

  detail(l): void {
    this.drawer.create({
      nzTitle: '详情',
      // nzWidth: 400,
      nzContent: LinkDetailComponent,
      nzContentParams: {
        link: l
      }
    });
  }
}
